package image

import (
	"net/url"
	"regexp"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alitto/pond"
	"github.com/sirupsen/logrus"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/file"
)

type Logic interface {
	Find(ctx *context.Context, filters url.Values) (files LocalImages, err error)
	Save(ctx *context.Context, files LocalImages, replace bool, batchSize int) error
	ReadFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, reBarcode *regexp.Regexp, filters url.Values) (files MapImages, err error)
	ReadBarcodeImageFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, filters url.Values) (files MapImages, err error)
	Sync(ctx *context.Context, root string, maxDepth int, replace bool, filters url.Values) (images MapImages, err error)
}

type logic struct {
	conn   *simutils.DBConnection
	client *client.Client
	repo   Repo
}

func newLogic(conn *simutils.DBConnection, client *client.Client, repo Repo) (Logic, error) {
	l := &logic{
		conn:   conn,
		client: client,
		repo:   repo,
	}
	return l, nil
}

func (l *logic) Find(ctx *context.Context, filters url.Values) (files LocalImages, err error) {
	filters["content_types"] = []string{ImageContentTypeRegex}
	q := l.conn.DB.WithContext(ctx.Request().Context())

	if storedFiles, err := l.repo.Read(ctx, q, filters); err != nil {
		return nil, err
	} else {
		return storedFiles.GetValues(), nil
	}
}

func (l *logic) Save(ctx *context.Context, images LocalImages, replace bool, batchSize int) error {
	pool := pond.New(batchSize, 0)

	for _, img := range images {
		if img.File == nil || img.File.Synced() {
			continue
		}

		pool.Submit(func() {
			logrus.Infof("Running task for %v", img)
			defer logrus.Infof("Finished task for %v", img)
			// Upload files
			if image, err := l.client.StoreImage(ctx,
				img.File.Path,
				img.Image,
			); err != nil {
				return
			} else if img.Image = image; image == nil {
				return
			} else if tx := l.conn.DB.WithContext(ctx.Request().Context()); tx == nil {
				return
			} else if err := l.repo.Create(ctx, tx, LocalImages{img}); err != nil {
				logrus.Infof("writing file %v in db failed", img)
				return
			}
		})
	}

	pool.StopAndWait()

	return nil
}

func (l *logic) ReadFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, reBarcode *regexp.Regexp, filters url.Values) (files MapImages, err error) {
	if reContentType, err := regexp.Compile(ImageContentTypeRegex); err != nil {
		return nil, err
	} else if reBarcode, err := regexp.Compile(ImageBarcodeRegex); err != nil {
		return nil, err
	} else if filteredImages, _ := file.WalkDir(root, maxDepth, reContentType, nil); len(filteredImages) == 0 {
		logrus.Info("!!! no images found !!!")
		return nil, nil
	} else if _, files = FromFiles(filteredImages, reBarcode); len(files) == 0 {
		logrus.Info("xxx convert File to Image failed xxx")
		return files, nil
	} else {
		return files, nil
	}
}

func (l *logic) ReadBarcodeImageFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, filters url.Values) (files MapImages, err error) {
	if reBarcode, err := regexp.Compile(ImageBarcodeRegex); err != nil {
		return nil, err
	} else if mapImages, err := l.ReadFiles(ctx, root, maxDepth, reContentType, reBarcode, filters); err != nil {
		return nil, nil
	} else if q := l.conn.DB.WithContext(ctx.Request().Context()); q == nil {
		return nil, simutils.ErrInvalidDatabaseConnection
	} else if _, err := l.repo.ReadFiles(ctx, q, mapImages, filters); err != nil {
		return nil, err
	} else {
		return mapImages, nil
	}
}

func (l *logic) Sync(ctx *context.Context, root string, maxDepth int, replace bool, filters url.Values) (images MapImages, err error) {
	if mfiles, err := l.ReadBarcodeImageFiles(ctx, root, maxDepth, nil, filters); err != nil {
		return nil, err
	} else if len(mfiles) == 0 {
		return nil, nil
	} else if err := l.Save(ctx, mfiles.GetValues(), replace, 5); err != nil {
		return nil, err
	} else {
		return mfiles, nil
	}
}
