package upgrades

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type UpgradeSourceConfig struct {
	Region         string `yaml:"region"`
	S3BucketName   string `yaml:"s3_bucket_name"`
	CommonPrefix   string `yaml:"common_prefix"`
	DestinationDir string `yaml:"destination_dir"`
}

type DownloadRequest struct {
	SourcePath string `json:"source_path"`

	// will get scoped under a managed path
	DestinationPath string `json:"destination_path"`

	// whether periodic cleanups should ignore the downloaded file
	DisableAutoCleanup bool `json:"disable_auto_cleanup"`
}

type DownloadResponse struct {
	// absolute path to where the file will be downloaded to
	DownloadFileDestination string `json:"download_file_destination"`
}

type Upgrades struct {
	Config *UpgradeSourceConfig
}

func (u *Upgrades) absoluteDestinationFromSource(requestedDestinationPath string, suppressCleanup bool) string {
	persistenceMode := "transient"
	if suppressCleanup == true {
		persistenceMode = "persistent"
	}
	return path.Join(u.Config.DestinationDir, persistenceMode, requestedDestinationPath)
}

func getFreeDiskSpace() (uint64, error) {
	var stat unix.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0, errors.Wrap(err, "could not get cwd")
	}
	err = unix.Statfs(wd, &stat)
	if err != nil {
		return 0, errors.Wrap(err, "failed to stat cwd")
	}
	// Available blocks * size per block = available space in bytes
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	return freeBytes, nil
}

func (u *Upgrades) DownloadFile(r *DownloadRequest) (*DownloadResponse, error) {
	destination := u.absoluteDestinationFromSource(r.DestinationPath, r.DisableAutoCleanup)
	if _, err := os.Stat(destination); err == nil {
		log.Printf("file %s already exists; removing\n", destination)
		err := os.Remove(destination)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to remove %s", destination)
		}
	}

	bucketName := aws.String(u.Config.S3BucketName)
	sourcePath := aws.String(path.Join(u.Config.CommonPrefix, r.SourcePath))

	// initiate download
	s3Client := s3.NewFromConfig(aws.Config{Region: u.Config.Region})
	ctx := context.Background()

	// check that we have enough space
	attrs, err := s3Client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
		ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesObjectSize},
		Bucket:           bucketName,
		Key:              sourcePath,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get attrs for %s/%s", *bucketName, *sourcePath)
	}

	freeSpace, err := getFreeDiskSpace()
	if err != nil {
		return nil, errors.Wrap(err, "could not determine free disk space")
	}

	if uint64(attrs.ObjectSize) > freeSpace {
		return nil, fmt.Errorf("request object %s/%s has size %d , while we only have free disk space %d bytes", *bucketName, *sourcePath, uint64(attrs.ObjectSize), freeSpace)
	}

	// Download the S3 object using the S3 manager object downloader
	fl, err := os.Create(destination)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s", destination)
	}
	downloader := manager.NewDownloader(s3Client)
	_, err = downloader.Download(context.TODO(), fl, &s3.GetObjectInput{
		Bucket: bucketName,
		Key:    sourcePath,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download %s/%s", *bucketName, *sourcePath)
	}
	return &DownloadResponse{
		DownloadFileDestination: destination,
	}, nil
}
