package e2etest

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/nefeli/terraform-exec/tfexec"
)

const defaultWorkspace = "default"

func TestWorkspace_default_only(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		assertWorkspaceList(t, tf, defaultWorkspace)
		assertWorkspaceShow(t, tf, defaultWorkspace)

		t.Run("select default when already on default", func(t *testing.T) {
			err := tf.WorkspaceSelect(context.Background(), defaultWorkspace)
			if err != nil {
				t.Fatalf("unable to select workspace: %s", err)
			}

			assertWorkspaceList(t, tf, defaultWorkspace)
			assertWorkspaceShow(t, tf, defaultWorkspace)
		})

		t.Run("create new workspace", func(t *testing.T) {
			const newWorkspace = "new1"
			err := tf.WorkspaceNew(context.Background(), newWorkspace)
			if err != nil {
				t.Fatalf("got error creating new workspace: %s", err)
			}

			assertWorkspaceList(t, tf, newWorkspace, newWorkspace)
			assertWorkspaceShow(t, tf, newWorkspace)
		})
	})
}

func TestWorkspace_does_not_exist(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		const doesNotExistWorkspace = "does-not-exist"
		err := tf.WorkspaceSelect(context.Background(), doesNotExistWorkspace)
		var wsErr *tfexec.ErrNoWorkspace
		if !errors.As(err, &wsErr) {
			t.Fatalf("expected ErrNoWorkspace, %T returned: %s", err, err)
		}
		if wsErr.Name != doesNotExistWorkspace {
			t.Fatalf("expected %q, got %q", doesNotExistWorkspace, wsErr.Name)
		}
	})
}

func TestWorkspace_already_exists(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		const newWorkspace = "existing-workspace"
		t.Run("create new workspace", func(t *testing.T) {
			err := tf.WorkspaceNew(context.Background(), newWorkspace)
			if err != nil {
				t.Fatalf("got error creating new workspace: %s", err)
			}

			assertWorkspaceList(t, tf, newWorkspace, newWorkspace)
			assertWorkspaceShow(t, tf, newWorkspace)
		})

		t.Run("create existing workspace", func(t *testing.T) {
			err := tf.WorkspaceNew(context.Background(), newWorkspace)

			var wsErr *tfexec.ErrWorkspaceExists
			if !errors.As(err, &wsErr) {
				t.Fatalf("expected ErrWorkspaceExists, %T returned: %s", err, err)
			}
			if wsErr.Name != newWorkspace {
				t.Fatalf("expected %q, got %q", newWorkspace, wsErr.Name)
			}
		})
	})
}

func TestWorkspace_multiple(t *testing.T) {
	runTest(t, "workspaces", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		assertWorkspaceList(t, tf, "foo", "foo")
		assertWorkspaceShow(t, tf, "foo")

		const newWorkspace = "new1"

		t.Run("create new workspace", func(t *testing.T) {
			err := tf.WorkspaceNew(context.Background(), newWorkspace)
			if err != nil {
				t.Fatalf("got error creating new workspace: %s", err)
			}

			assertWorkspaceList(t, tf, newWorkspace, "foo", newWorkspace)
			assertWorkspaceShow(t, tf, newWorkspace)
		})

		t.Run("select non-default workspace", func(t *testing.T) {
			err := tf.WorkspaceSelect(context.Background(), "foo")
			if err != nil {
				t.Fatalf("unable to select workspace: %s", err)
			}

			assertWorkspaceList(t, tf, "foo", "foo", newWorkspace)
			assertWorkspaceShow(t, tf, "foo")
		})
	})
}

func TestWorkspace_deletion(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		assertWorkspaceList(t, tf, defaultWorkspace)
		assertWorkspaceShow(t, tf, defaultWorkspace)

		const testWorkspace = "testws"

		t.Run("create and delete workspace", func(t *testing.T) {
			err := tf.WorkspaceNew(context.Background(), testWorkspace)
			if err != nil {
				t.Fatalf("got error creating workspace: %s", err)
			}

			assertWorkspaceList(t, tf, testWorkspace, testWorkspace)
			assertWorkspaceShow(t, tf, testWorkspace)

			err = tf.WorkspaceSelect(context.Background(), defaultWorkspace)
			if err != nil {
				t.Fatalf("got error selecting workspace: %s", err)
			}

			assertWorkspaceShow(t, tf, defaultWorkspace)

			err = tf.WorkspaceDelete(context.Background(), testWorkspace)
			if err != nil {
				t.Fatalf("got error deleting workspace: %s", err)
			}

			assertWorkspaceList(t, tf, defaultWorkspace)
			assertWorkspaceShow(t, tf, defaultWorkspace)
		})
	})
}

func assertWorkspaceList(t *testing.T, tf *tfexec.Terraform, expectedCurrent string, expectedWorkspaces ...string) {
	actualWorkspaces, actualCurrent, err := tf.WorkspaceList(context.Background())
	if err != nil {
		t.Fatalf("got error querying workspace list: %s", err)
	}
	if actualCurrent != expectedCurrent {
		t.Fatalf("expected %q workspace to be selected, got %q", expectedCurrent, actualCurrent)
	}
	expectedWorkspaces = append([]string{defaultWorkspace}, expectedWorkspaces...)
	if !reflect.DeepEqual(actualWorkspaces, expectedWorkspaces) {
		t.Fatalf("expected %#v, got %#v", actualWorkspaces, expectedWorkspaces)
	}
}

func assertWorkspaceShow(t *testing.T, tf *tfexec.Terraform, expectedWorkspace string) {
	actualWorkspace, err := tf.WorkspaceShow(context.Background())
	if err != nil {
		t.Fatalf("got error querying workspace show: %s", err)
	}
	if actualWorkspace != expectedWorkspace {
		t.Fatalf("expected %q workspace, got %q workspace", expectedWorkspace, actualWorkspace)
	}
}
