package drivers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/scilive/scibase/drivers"
	"github.com/stretchr/testify/assert"
)

// https://docs.storj.io/dcs/how-tos/host-a-static-website/host-a-static-website-with-the-cli-and-linksharing-service
func TestNewMinIO(t *testing.T) {
	s3Client, err := drivers.NewMinIO()
	assert.Nil(t, err)
	buckets, err := s3Client.ListBuckets(context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, buckets)
	fmt.Println(buckets)

}
