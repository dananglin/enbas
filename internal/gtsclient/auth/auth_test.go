package auth_test

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient/auth"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	t.Log("Loading the test authentication details from file.")

	path := filepath.Join(
		"testdata",
		t.Name()+".golden",
	)

	testAuth, err := auth.NewAuthFromConfig(path)
	if err != nil {
		t.Fatalf(
			"FAILED test %q: Error loading the test authentication from %s: %v",
			t.Name(),
			path,
			err,
		)
	}

	compareAuths(t, testAuth, authCompare{
		instanceURL:      "https://gts-01.social.example",
		token:            "WGEUCBJXDUSQULCX6CT7244EVFUQQT3SVLDQKCFAGWII01MY",
		currentAccountID: "",
	})

	t.Log("Updating the current account ID.")

	newCurrentAccontID := "01VRM0GUDCBXCMWMJZJBCAGB0IGYDYRQ"

	// Run the update twice for full coverage
	testAuth.UpdateCurrentAccountID(newCurrentAccontID)
	testAuth.UpdateCurrentAccountID(newCurrentAccontID)

	compareAuths(t, testAuth, authCompare{
		instanceURL:      "https://gts-01.social.example",
		token:            "WGEUCBJXDUSQULCX6CT7244EVFUQQT3SVLDQKCFAGWII01MY",
		currentAccountID: newCurrentAccontID,
	})

	t.Log("Update the authentication details with new values")

	newAuthCfg := config.Credentials{
		Instance:     "https://gts-03.social.example",
		ClientID:     "01RWQWPO9PYKDYCLSGVA4AKZWJ6T8VPI",
		ClientSecret: "88d473b1-3f1f-4ae2-ae85-5e5fa30d1b77",
		AccessToken:  "NLYDVCERZHSCFBIPNE5OOV9KBJXS3RFLC8PRC9ABXXXZSJFA",
	}

	testAuth.UpdateAuth(newAuthCfg)

	compareAuths(t, testAuth, authCompare{
		instanceURL:      "https://gts-03.social.example",
		token:            "NLYDVCERZHSCFBIPNE5OOV9KBJXS3RFLC8PRC9ABXXXZSJFA",
		currentAccountID: "",
	})
}

func TestNewAuthZero(t *testing.T) {
	t.Parallel()

	testAuth := auth.NewAuthZero()

	compareAuths(t, testAuth, authCompare{
		instanceURL:      "",
		token:            "",
		currentAccountID: "",
	})
}

func TestNewAuthWithNoCurrentAccount(t *testing.T) {
	t.Parallel()

	path := filepath.Join(
		"testdata",
		t.Name()+".golden",
	)

	testAuth, err := auth.NewAuthFromConfig(path)
	if err != nil {
		t.Fatalf(
			"FAILED test %q: Error loading the test authentication from %s: %v",
			t.Name(),
			path,
			err,
		)
	}

	compareAuths(t, testAuth, authCompare{
		instanceURL:      "",
		token:            "",
		currentAccountID: "",
	})
}

func TestNoAuthenticationDetailsErrorHandling(t *testing.T) {
	t.Parallel()

	path := filepath.Join(
		"testdata",
		t.Name()+".golden",
	)

	_, err := auth.NewAuthFromConfig(path)
	if err == nil {
		t.Errorf(
			"FAILED test %q: No error received after loading the test authentication with an incorrect current account.",
			t.Name(),
		)

		return
	}

	wantErr := auth.NewNoAuthenticationDetailsError("fake-account@gts.fake.example")

	if !errors.Is(err, wantErr) {
		t.Errorf("FAILED test %q: Unexpected error received after loading the test authentication with an incorrect current account.\nwant %v\n got: %v",
			t.Name(),
			wantErr,
			err,
		)
	} else {
		t.Logf("GOOD result from %q: Expected error received after loading the test authentication with an incorrect current account.\ngot: %v",
			t.Name(),
			err,
		)
	}
}

type authCompare struct {
	instanceURL      string
	token            string
	currentAccountID string
}

func compareAuths(t *testing.T, gotAuth *auth.Auth, want authCompare) {
	t.Helper()

	got := authCompare{
		instanceURL:      gotAuth.GetInstanceURL(),
		token:            gotAuth.GetToken(),
		currentAccountID: gotAuth.GetCurrentAccountID(),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"BAD results from \"compareAuths\" for %q: Unexpected values found after comparing the test auth.\nwant: %+v\n got: %+v",
			t.Name(),
			want,
			got,
		)
	} else {
		t.Logf(
			"GOOD result from \"compareAuths\" for %q: Expected values found after comparing the test auth.\ngot: %+v",
			t.Name(),
			got,
		)
	}
}
