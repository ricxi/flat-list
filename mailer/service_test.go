package mailer

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockMailer struct {
	out string
	err error
}

func (m *mockMailer) send(from, to, subject, body string) error {
	m.out = body
	return m.err
}

// I haven't really written any cases where the
// mailer returns an error yet.
func TestServiceSendActivationEmail(t *testing.T) {
	type args struct {
		data ActivationEmailData
	}

	type expected struct {
		errString string
	}

	testCases := []struct {
		name     string
		service  Service
		args     args
		expected expected
	}{
		{
			name: "Success",
			service: Service{
				mailer: &mockMailer{
					err: nil,
				},
				emailTemplatesDir: "./templates",
			},
			args: args{
				data: ActivationEmailData{
					To:      "michaelscott@dundermifflin.com",
					From:    "theteam@flatlist.com",
					Subject: "How to activate your new account",
					ActivationData: ActivationData{
						Name:      "Michael",
						Hyperlink: "http://localhost:5000/clickme",
					},
				},
			},
			expected: expected{
				errString: "",
			},
		},
		{
			name:    "MissingToField",
			service: Service{mailer: nil, emailTemplatesDir: ""},
			args: args{
				data: ActivationEmailData{
					From:    "theteam@flatlist.com",
					Subject: "How to activate your new account",
					ActivationData: ActivationData{
						Name:      "Michael",
						Hyperlink: "http://localhost:5000/clickme",
					},
				},
			},
			expected: expected{
				errString: "missing field is required: to",
			},
		},
		{
			name:    "MissingFromField",
			service: Service{mailer: nil, emailTemplatesDir: ""},
			args: args{
				data: ActivationEmailData{
					To:      "michaelscott@dundermifflin.com",
					Subject: "How to activate your new account",
					ActivationData: ActivationData{
						Name:      "Michael",
						Hyperlink: "http://localhost:5000/clickme",
					},
				},
			},
			expected: expected{
				errString: "missing field is required: from",
			},
		},
		{
			name:    "MissingSubjectField",
			service: Service{mailer: nil, emailTemplatesDir: ""},
			args: args{
				data: ActivationEmailData{
					From: "theteam@flatlist.com",
					To:   "michaelscott@dundermifflin.com",
					ActivationData: ActivationData{
						Name:      "Michael",
						Hyperlink: "http://localhost:5000/clickme",
					},
				},
			},
			expected: expected{
				errString: "missing field is required: subject",
			},
		},
		{
			name:    "MissingActivationHyperLinkField",
			service: Service{mailer: nil, emailTemplatesDir: ""},
			args: args{
				data: ActivationEmailData{
					From:    "theteam@flatlist.com",
					To:      "michaelscott@dundermifflin.com",
					Subject: "How to activate your new account",
					ActivationData: ActivationData{
						Name: "Michael",
					},
				},
			},
			expected: expected{
				errString: "missing field is required: activationHyperlink",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.sendActivationEmail(tt.args.data)
			if err != nil {
				assert.NotEmpty(t, err)
				assert.EqualError(t, err, tt.expected.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// This is also a test for the Service's sendActivationEmail
// method, but it checks the contents of the email generated
func TestServiceSendActivationEmailContents(t *testing.T) {
	t.Run("SendWithName", func(t *testing.T) {
		const (
			inputEmailTmpl = "./testdata"
			goldenFile     = "./testdata/useractivation.goldenfile.html"
		)

		// mockMailerDest is used to get the email body
		mockMailerDst := mockMailer{}
		service := &Service{
			emailTemplatesDir: inputEmailTmpl,
			mailer:            &mockMailerDst,
		}

		data := ActivationEmailData{
			To:      "michaelscott@dundermifflin.com",
			From:    "theteam@flatlist.com",
			Subject: "How to activate your new account",
			ActivationData: ActivationData{
				Name:      "Michael",
				Hyperlink: "http://localhost:5000/clickme",
			},
		}

		err := service.sendActivationEmail(data)
		assert.NoError(t, err)

		expected, err := os.ReadFile(goldenFile)
		require.NoError(t, err)

		assert.Equal(t, expected, []byte(mockMailerDst.out))
		// uncomment to get more descriptive messages for failing cases
		// assert.Equal(t, string(expected), mockMailerDst.out)
		// assert.Equal(t, expected, []byte(mockMailerDst.out))
	})

	t.Run("SendWithoutName", func(t *testing.T) {
		const (
			inputEmailTmpl = "./testdata"
			goldenFile     = "./testdata/useractivation.noname.goldenfile.html"
		)

		mockMailerDst := mockMailer{}
		service := &Service{
			emailTemplatesDir: inputEmailTmpl,
			mailer:            &mockMailerDst,
		}

		data := ActivationEmailData{
			To:      "michaelscott@dundermifflin.com",
			From:    "theteam@flatlist.com",
			Subject: "How to activate your new account",
			ActivationData: ActivationData{
				Name:      "", // the Name field is left intentionally blank
				Hyperlink: "http://localhost:5000/clickme",
			},
		}

		err := service.sendActivationEmail(data)
		assert.NoError(t, err)

		expected, err := os.ReadFile(goldenFile)
		require.NoError(t, err)

		assert.True(t, bytes.Equal(expected, []byte(mockMailerDst.out)))
		// uncomment to get more descriptive messages for failing cases
		// assert.Equal(t, string(expected), mockMailerDst.out)
		// assert.Equal(t, expected, []byte(mockMailerDst.out))
	})
}
