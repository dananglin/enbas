package sessions_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/sessions"
)

func TestInitialisedSessionStore(t *testing.T) {
	sessionSet := sessions.NewSessionStore(true)

	var initialised bool
	if err := sessionSet.IsInitialised(
		sessions.NoRPCArgs{},
		&initialised,
	); err != nil {
		t.Fatalf(
			"FAILED test %s: Received an unexpected error after checking to see if the SessionStore is initialised: %v",
			t.Name(),
			err,
		)
	}

	if !initialised {
		t.Fatalf(
			"FAILED test %s: Unexpected result received after checking to see if the SessionStore is initialised.\nwant: %t\n got: %t",
			t.Name(),
			true,
			initialised,
		)
	}

	if !sessionSet.IsEmpty() {
		t.Fatalf(
			"FAILED test %s: unexpected result received after checking to see if the initialised SessionStore is empty.\nwant: %t\n got: %t",
			t.Name(),
			true,
			sessionSet.IsEmpty(),
		)
	}

	t.Log("Successfully initialised the SessionStore.")

	// Concurrently add 20 session IDs to the SessionStore
	sessionIDs := []string{
		"B6ndL6V8FXr99vHorfzcCjy7tmHhbMqK", "wiZV8Gn1chZei2DoYLGkiYiyUeDAEqX2",
		"T8V0HzXix3w3U1NuNXDra9YHN4Tr5M77", "Jbx2ZLpVDfwHnHXyQGq2OCoE9nia4beL",
		"M61qkaxo6zjC2Hb2A9t3x6zUZfx99Nyv", "VoUvODNXuiWHIGpwvJw8fjEQtvzTvR2j",
		"IEZ7X8Lkjn303PtXge0jFNOFd1xrGkBH", "UgloUqn3icS4OAEP4bhHGHKRVS4wL59x",
		"poW7FfIACsIM9USoLwov2B3YPtUOOX55", "CFlc9NU8ZpuL4mOLMQ1csMcjekEE9Kcl",
		"AiWNFCPXw8oCYwpWe3Jdobn57ZM8k2cR", "KbwoVZ0t4yUtj4LNYXvq4ZIwqBsBcdQn",
		"YQqDFpGLeX4LkSc4cPnQfteXZR2yh7Rm", "A1vODi1w0hBubYfKFBWqsrIXkRb6IPYu",
		"KsIoYeoJnxJDKeXwGbK5E8qyqkqpRZuO", "9syzSzSXdadeYCeCc1ChDreenTT48gTB",
		"MqWK6jQFsPFuL74nPiacXr74lVc3lyzY", "FMNQqCa2D0wtT4sy97oS7RGqgRqo9jvb",
		"jrB3vSTAYW7hEB7svMZOzybDXNrNYxRr", "mcVWXylY5qaDdiMHCeJCAp3G91IWrF1q",
	}

	var (
		errs []error
		mu   sync.Mutex
		wg   sync.WaitGroup
	)

	wg.Add(len(sessionIDs))

	for idx := range sessionIDs {
		go func() {
			defer wg.Done()

			err := sessionSet.Add(sessionIDs[idx], nil)
			if err != nil {
				mu.Lock()
				{
					errs = append(errs, err)
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		msg := fmt.Sprintf("FAILED test %s: Received one or more errors after adding the session IDs to the SessionStore.\n", t.Name())

		for idx := range errs {
			msg += fmt.Sprintf("Error %d: %s\n", idx+1, errs[idx].Error())
		}

		t.Fatal(msg)
	}

	t.Log("Successfully added the session IDs to the SessionStore")

	if sessionSet.NumSessionIDs() != len(sessionIDs) {
		t.Fatalf(
			"FAILED test %s: Unexpected number of session IDs found in the SessionStore.\nwant: %d\n got: %d",
			t.Name(),
			len(sessionIDs),
			sessionSet.NumSessionIDs(),
		)
	} else {
		t.Logf(
			"GOOD result from %s: Expected number of session IDs found in the SessionStore.\ngot: %d",
			t.Name(),
			sessionSet.NumSessionIDs(),
		)
	}

	// Ensure that the SessionStore prevents attempts to
	// add duplicate session IDs to the SessionStore.
	err := sessionSet.Add(sessionIDs[0], nil)
	if err == nil {
		t.Errorf(
			"FAILED test %s: No error received after adding a duplicate session ID to the SessionStore.",
			t.Name(),
		)
	}

	wantErr := sessions.NewSessionIDExistsError(sessionIDs[0])

	if !errors.Is(err, wantErr) {
		t.Errorf(
			"FAILED test %s: Unexpected error received after attempting to add a duplicate session ID to the SessionStore.\nwant: %v\n got: %v",
			t.Name(),
			wantErr,
			err,
		)
	} else {
		t.Logf(
			"GOOD result from %s: Expected error received after attempting to add a duplicate session ID to the SessionStore.\ngot: %v",
			t.Name(),
			err,
		)
	}

	// Concurrently remove the session IDs to the SessionStore
	wg.Add(len(sessionIDs))

	for idx := range sessionIDs {
		go func() {
			defer wg.Done()

			err := sessionSet.Remove(sessionIDs[idx], nil)
			if err != nil {
				mu.Lock()
				{
					errs = append(errs, err)
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		msg := fmt.Sprintf("FAILED test %s: Received one or more errors after removing the session IDs to the SessionStore.\n", t.Name())

		for idx := range errs {
			msg += fmt.Sprintf("Error %d: %s\n", idx+1, errs[idx].Error())
		}

		t.Fatal(msg)
	}

	t.Log("Successfully removed the session IDs to the SessionStore")

	if !sessionSet.IsEmpty() {
		t.Errorf(
			"FAILED test %s: The SessionStore does not appear to be empty after removing all the session IDs.\nNumber of session IDs remaining: %d",
			t.Name(),
			sessionSet.NumSessionIDs(),
		)
	} else {
		t.Logf(
			"GOOD result from %s: The SessionStore is empty after removing all session IDs.\nNumber of session IDs remaining: %d",
			t.Name(),
			sessionSet.NumSessionIDs(),
		)
	}
}

func TestNonInitialisedSessionStore(t *testing.T) {
	sessionSet := sessions.NewSessionStore(false)

	var initialised bool
	if err := sessionSet.IsInitialised(
		sessions.NoRPCArgs{},
		&initialised,
	); err != nil {
		t.Fatalf(
			"FAILED test %s: Received an unexpected error after checking to see if the SessionStore is initialised: %v",
			t.Name(),
			err,
		)
	}

	if initialised {
		t.Fatalf(
			"FAILED test %s: Unexpected result received after checking to see if the non-initialised SessionStore is initialised.\nwant: %t\n got: %t",
			t.Name(),
			false,
			initialised,
		)
	}

	t.Log("Successfully created the non-initialised the SessionStore.")

	// Ensure that the non-initialised SessionStore prevents any
	// attempts to add or remove session IDs.
	wantErr := sessions.NonInitialisedSessionStoreError{}

	err := sessionSet.Add("dgInuwqLBm4GlrN1jrjEx09kgR2CkF1b", nil)
	if err == nil {
		t.Fatalf(
			"FAILED test %s: No error received after attempting to add a session ID to a non-initialised SessionStore",
			t.Name(),
		)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"FAILED test %s: Unexpected error received after attempting to add a session ID to a non-initialised SessionStore.\nwant: %v\n got: %v",
			t.Name(),
			wantErr,
			err,
		)
	}

	t.Logf(
		"GOOD result from %s: Expected error received after attempting to add a session ID to a non-initialised SessionStore.\ngot: %v",
		t.Name(),
		err,
	)

	err = sessionSet.Remove("dgInuwqLBm4GlrN1jrjEx09kgR2CkF1b", nil)
	if err == nil {
		t.Fatalf(
			"FAILED test %s: No error received after attempting to remove a session ID from a non-initialised SessionStore",
			t.Name(),
		)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"FAILED test %s: Unexpected error received after attempting to remove a session ID from a non-initialised SessionStore.\nwant: %v\n got: %v",
			t.Name(),
			wantErr,
			err,
		)
	}

	t.Logf(
		"GOOD result from %s: Expected error received after attempting to remove a session ID from a non-initialised SessionStore.\ngot: %v",
		t.Name(),
		err,
	)
}
