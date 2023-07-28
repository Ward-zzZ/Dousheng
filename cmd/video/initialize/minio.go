package initialize

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"tiktok-demo/cmd/video/config"
)

func InitMinio() *minio.Client {
	mi := config.GlobalServerConfig.MinioInfo
	// Initialize minio client object.
	// endpoint : 对象存储服务的URL
	// accessKeyID : 标识的用户唯一ID
	// secretAccessKey : 访问密钥
	// Bucket : 存储桶名称
	klog.Infof("minio url: %s", "127.0.0.1:9000")
	mc, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(mi.AccessKeyID, mi.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		klog.Fatalf("create minio client err: %s", err.Error())
	}
	bucketName := []string{config.GlobalServerConfig.MinioInfo.VideoBucket, config.GlobalServerConfig.MinioInfo.CoverBucket}
	for _, bucket := range bucketName {
		exists, err := mc.BucketExists(context.Background(), bucket)
		if err != nil {
			klog.Fatal("check bucket exists err: %s", err.Error())
		}
		if !exists {
			err = mc.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{Region: "cn-north-1"})
			if err != nil {
				klog.Fatalf("make bucket err: %s", err.Error())
			}
		}
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":["arn:aws:s3:::` + bucket + `"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucket + `/*"]}]}`
		err = mc.SetBucketPolicy(context.Background(), bucket, policy)
		if err != nil {
			klog.Fatal("set bucket policy err:%s", err)
		}
	}
	return mc
}
