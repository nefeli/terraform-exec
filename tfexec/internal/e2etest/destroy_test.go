package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/nefeli/terraform-exec/tfexec"
)

func TestDestroy(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		err = tf.Destroy(context.Background())
		if err != nil {
			t.Fatalf("error running Destroy: %s", err)
		}
	})
}
