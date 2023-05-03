package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/app/mocks"
	"homework10/internal/users"
)

type HTTPGINTestSuite struct {
	suite.Suite
	a          *mocks.App
	r          *httptest.ResponseRecorder
	c          *gin.Context
	setReqBody func(string, map[string]any)
}

func (s *HTTPGINTestSuite) SetupSuite() {
	s.a = mocks.NewApp(s.T())
	gin.SetMode(gin.TestMode)

	s.setReqBody = func(meth string, m map[string]any) {
		data, _ := json.Marshal(m)
		s.c.Request = httptest.NewRequest(meth, "http://not.nil.url", bytes.NewReader(data))
		s.c.Request.Header.Add("Content-Type", "application/json")
	}
}

func (s *HTTPGINTestSuite) SetupSubTest() {
	s.r = httptest.NewRecorder()
	c, _ := gin.CreateTestContext(s.r)
	s.c = c
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_CreateAd() {
	handler := createAd(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		reqBody map[string]any
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "title",
				"text":    "text",
			},
			setMock: func() {
				s.a.
					On("CreateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "title",
				"text":    "text",
			},
			setMock: func() {
				s.a.
					On("CreateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "title",
				"text":    "text",
			},
			setMock: func() {
				s.a.
					On("CreateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "title",
				"text":    "text",
			},
			setMock: func() {
				s.a.
					On("CreateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(&ads.Ad{
						ID:        0,
						Title:     "title",
						Text:      "text",
						UserID:    0,
						Published: false,
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": adResponse{
						ID:        0,
						Title:     "title",
						Text:      "text",
						AuthorID:  0,
						Published: false,
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.setReqBody(http.MethodPost, tt.reqBody)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_ChangeAdStatus() {
	handler := changeAdStatus(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		reqBody map[string]any
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			reqBody: map[string]any{
				"user_id":   0,
				"published": true,
			},
			setMock: func() {
				s.a.
					On("ChangeAdStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			reqBody: map[string]any{
				"user_id":   0,
				"published": true,
			},
			setMock: func() {
				s.a.
					On("ChangeAdStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			reqBody: map[string]any{
				"user_id":   0,
				"published": true,
			},
			setMock: func() {
				s.a.
					On("ChangeAdStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			reqBody: map[string]any{
				"user_id":   0,
				"published": true,
			},
			setMock: func() {
				s.a.
					On("ChangeAdStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(&ads.Ad{
						ID:        0,
						Title:     "title",
						Text:      "text",
						UserID:    0,
						Published: true,
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": adResponse{
						ID:        0,
						Title:     "title",
						Text:      "text",
						AuthorID:  0,
						Published: true,
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("ad_id", "0")
			s.setReqBody(http.MethodPut, tt.reqBody)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_UpdateAd() {
	handler := updateAd(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		reqBody map[string]any
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "new title",
				"text":    "new text",
			},
			setMock: func() {
				s.a.
					On("UpdateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "new title",
				"text":    "new text",
			},
			setMock: func() {
				s.a.
					On("UpdateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "new title",
				"text":    "new text",
			},
			setMock: func() {
				s.a.
					On("UpdateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			reqBody: map[string]any{
				"user_id": 0,
				"title":   "new title",
				"text":    "new text",
			},
			setMock: func() {
				s.a.
					On("UpdateAd", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(&ads.Ad{
						ID:        0,
						Title:     "title",
						Text:      "text",
						UserID:    0,
						Published: true,
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": adResponse{
						ID:        0,
						Title:     "title",
						Text:      "text",
						AuthorID:  0,
						Published: true,
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("ad_id", "0")
			s.setReqBody(http.MethodPut, tt.reqBody)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_ShowAd() {
	handler := showAd(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			setMock: func() {
				s.a.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			setMock: func() {
				s.a.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			setMock: func() {
				s.a.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			setMock: func() {
				s.a.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{
						ID:        0,
						Title:     "title",
						Text:      "text",
						UserID:    0,
						Published: true,
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": adResponse{
						ID:        0,
						Title:     "title",
						Text:      "text",
						AuthorID:  0,
						Published: true,
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("ad_id", "0")
			s.setReqBody(http.MethodGet, nil)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_ListAds() {
	handler := listAds(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			setMock: func() {
				s.a.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			setMock: func() {
				s.a.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			setMock: func() {
				s.a.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			setMock: func() {
				s.a.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return([]*ads.Ad{
						{
							ID:        0,
							Title:     "1st title",
							Text:      "1st text",
							UserID:    0,
							Published: true,
						},
						{
							ID:        1,
							Title:     "2nd title",
							Text:      "2nd text",
							UserID:    0,
							Published: true,
						},
						{
							ID:        2,
							Title:     "3rd title",
							Text:      "3rd text",
							UserID:    0,
							Published: true,
						},
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": []adResponse{
						{
							ID:        0,
							Title:     "1st title",
							Text:      "1st text",
							AuthorID:  0,
							Published: true,
						},
						{
							ID:        1,
							Title:     "2nd title",
							Text:      "2nd text",
							AuthorID:  0,
							Published: true,
						},
						{
							ID:        2,
							Title:     "3rd title",
							Text:      "3rd text",
							AuthorID:  0,
							Published: true,
						},
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.setReqBody(http.MethodGet, nil)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_DeleteAd() {
	handler := deleteAd(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			setMock: func() {
				s.a.
					On("DeleteAd", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			setMock: func() {
				s.a.
					On("DeleteAd", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			setMock: func() {
				s.a.
					On("DeleteAd", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			setMock: func() {
				s.a.
					On("DeleteAd", mock.Anything, mock.Anything, mock.Anything).
					Return(&ads.Ad{
						ID:        0,
						Title:     "title",
						Text:      "text",
						UserID:    0,
						Published: true,
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": adResponse{
						ID:        0,
						Title:     "title",
						Text:      "text",
						AuthorID:  0,
						Published: true,
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("ad_id", "0")
			s.setReqBody(http.MethodDelete, nil)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_CreateUser() {
	handler := createUser(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		reqBody map[string]any
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			reqBody: map[string]any{
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			reqBody: map[string]any{
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			reqBody: map[string]any{
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			reqBody: map[string]any{
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(&users.User{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": userResponse{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.setReqBody(http.MethodPost, tt.reqBody)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_UpdateUser() {
	handler := updateUser(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		reqBody map[string]any
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			reqBody: map[string]any{
				"user_id":  0,
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			reqBody: map[string]any{
				"user_id":  0,
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			reqBody: map[string]any{
				"user_id":  0,
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			reqBody: map[string]any{
				"user_id":  0,
				"nickname": "user",
				"email":    "user@gmail.com",
			},
			setMock: func() {
				s.a.
					On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(&users.User{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": userResponse{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("user_id", "0")
			s.setReqBody(http.MethodPost, tt.reqBody)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_ShowUser() {
	handler := showUser(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			setMock: func() {
				s.a.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			setMock: func() {
				s.a.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			setMock: func() {
				s.a.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			setMock: func() {
				s.a.
					On("UserByID", mock.Anything, mock.Anything).
					Return(&users.User{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": userResponse{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("user_id", "0")
			s.setReqBody(http.MethodGet, nil)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func (s *HTTPGINTestSuite) TestHTTPGINHandlers_DeleteUser() {
	handler := deleteUser(s.a)

	type want struct {
		code int
		resp gin.H
	}
	tests := []struct {
		name    string
		setMock func()
		want    want
	}{
		{
			name: "forbidden error",
			setMock: func() {
				s.a.
					On("DeleteUser", mock.Anything, mock.Anything).
					Return(nil, app.ErrForbidden).
					Once()
			},
			want: want{
				code: http.StatusForbidden,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrForbidden.Error(),
				},
			},
		},
		{
			name: "bad request error",
			setMock: func() {
				s.a.
					On("DeleteUser", mock.Anything, mock.Anything).
					Return(nil, app.ErrBadRequest).
					Once()
			},
			want: want{
				code: http.StatusBadRequest,
				resp: gin.H{
					"data":  nil,
					"error": app.ErrBadRequest.Error(),
				},
			},
		},
		{
			name: "internal server error",
			setMock: func() {
				s.a.
					On("DeleteUser", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("untracked internal server error")).
					Once()
			},
			want: want{
				code: http.StatusInternalServerError,
				resp: gin.H{
					"data":  nil,
					"error": fmt.Errorf("untracked internal server error").Error(),
				},
			},
		},
		{
			name: "ok",
			setMock: func() {
				s.a.
					On("DeleteUser", mock.Anything, mock.Anything).
					Return(&users.User{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					}, nil).
					Once()
			},
			want: want{
				code: http.StatusOK,
				resp: gin.H{
					"data": userResponse{
						ID:       0,
						Nickname: "user",
						Email:    "user@gmail.com",
					},
					"error": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setMock()
			s.c.AddParam("user_id", "0")
			s.setReqBody(http.MethodDelete, nil)
			handler(s.c)
			data, _ := json.Marshal(tt.want.resp)
			assert.Equal(s.T(), tt.want.code, s.r.Code)
			assert.Equal(s.T(), data, s.r.Body.Bytes())
		})
	}
}

func TestHTTPGINTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPGINTestSuite))
}
