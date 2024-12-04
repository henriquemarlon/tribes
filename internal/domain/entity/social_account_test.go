package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSocialAccount(t *testing.T) {
	userID := uint(1)
	username := "test_user"
	followers := uint(100)
	platform := "twitter"
	createdAt := time.Now().Unix()

	tests := []struct {
		name             string
		userID           uint
		username         string
		followers        uint
		platform         string
		createdAt        int64
		expectedError    error
		errorMessagePart string
	}{
		{
			name:          "Valid Social Account",
			userID:        userID,
			username:      username,
			followers:     followers,
			platform:      platform,
			createdAt:     createdAt,
			expectedError: nil,
		},
		{
			name:             "Invalid User ID",
			userID:           0,
			username:         username,
			followers:        followers,
			platform:         platform,
			createdAt:        createdAt,
			expectedError:    ErrInvalidSocialAccount,
			errorMessagePart: "user ID cannot be zero",
		},
		{
			name:             "Empty Username",
			userID:           userID,
			username:         "",
			followers:        followers,
			platform:         platform,
			createdAt:        createdAt,
			expectedError:    ErrInvalidSocialAccount,
			errorMessagePart: "username cannot be empty",
		},
		{
			name:             "Zero Followers",
			userID:           userID,
			username:         username,
			followers:        0,
			platform:         platform,
			createdAt:        createdAt,
			expectedError:    ErrInvalidSocialAccount,
			errorMessagePart: "followers cannot be zero",
		},
		{
			name:             "Empty Platform",
			userID:           userID,
			username:         username,
			followers:        followers,
			platform:         "",
			createdAt:        createdAt,
			expectedError:    ErrInvalidSocialAccount,
			errorMessagePart: "platform cannot be empty",
		},
		{
			name:             "Invalid Platform",
			userID:           userID,
			username:         username,
			followers:        followers,
			platform:         "facebook",
			createdAt:        createdAt,
			expectedError:    ErrInvalidSocialAccount,
			errorMessagePart: "platform must be 'twitter' or 'instagram'",
		},
		{
			name:             "Missing Creation Date",
			userID:           userID,
			username:         username,
			followers:        followers,
			platform:         platform,
			createdAt:        0,
			expectedError:    ErrInvalidSocialAccount,
			errorMessagePart: "creation date is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socialAccount, err := NewSocialAccount(tt.userID, tt.username, tt.followers, tt.platform, tt.createdAt)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Contains(t, err.Error(), tt.errorMessagePart)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, socialAccount)
			}
		})
	}
}

func TestSocialAccount_Validate(t *testing.T) {
	validSocialAccount := &SocialAccount{
		UserId:    1,
		Username:  "valid_user",
		Followers: 100,
		Platform:  PlatformTwitter,
		CreatedAt: time.Now().Unix(),
	}

	t.Run("Valid Social Account", func(t *testing.T) {
		err := validSocialAccount.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		invalidAccount := *validSocialAccount
		invalidAccount.UserId = 0
		err := invalidAccount.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID cannot be zero")
	})

	t.Run("Empty Username", func(t *testing.T) {
		invalidAccount := *validSocialAccount
		invalidAccount.Username = ""
		err := invalidAccount.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username cannot be empty")
	})

	t.Run("Zero Followers", func(t *testing.T) {
		invalidAccount := *validSocialAccount
		invalidAccount.Followers = 0
		err := invalidAccount.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "followers cannot be zero")
	})

	t.Run("Invalid Platform", func(t *testing.T) {
		invalidAccount := *validSocialAccount
		invalidAccount.Platform = "facebook"
		err := invalidAccount.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "platform must be 'twitter' or 'instagram'")
	})
}
