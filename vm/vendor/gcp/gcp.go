package gcp

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// printImagesList prints a list of all non-deprecated image names available in given project.
func printImagesList(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %v", err)
	}
	defer imagesClient.Close()

	// Listing only non-deprecated images to reduce the size of the reply.
	req := &computepb.ListImagesRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(3),
		Filter:     proto.String("deprecated.state != DEPRECATED"),
	}

	// Although the `MaxResults` parameter is specified in the request, the iterator returned
	// by the `list()` method hides the pagination mechanic. The library makes multiple
	// requests to the API for you, so you can simply iterate over all the images.
	it := imagesClient.List(ctx, req)
	for {
		image, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s\n", image.GetName())
	}
	return nil
}
