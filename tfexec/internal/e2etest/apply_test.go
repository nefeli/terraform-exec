package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/nefeli/terraform-exec/tfexec"
)

var (
	validateMinApplyDestroyVersion = version.Must(version.NewVersion("0.15.2"))
)

func TestApply(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
	})
}

func TestApplyDestroy(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(validateMinApplyDestroyVersion) {
			t.Skip("terraform apply -destroy was added in Terraform 0.15.2, so test is not valid")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		err = tf.Apply(context.Background(), tfexec.Destroy(true))
		if err != nil {
			t.Fatalf("error running Apply -destroy: %s", err)
		}
	})
}
