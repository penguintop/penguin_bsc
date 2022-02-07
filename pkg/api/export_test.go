// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import "github.com/penguintop/penguin_bsc/pkg/penguin"

type Server = server

type (
	BytesPostResponse     = bytesPostResponse
	ChunkAddressResponse  = chunkAddressResponse
	SocPostResponse       = socPostResponse
	FeedReferenceResponse = feedReferenceResponse
	PenUploadResponse     = penUploadResponse
	TagResponse           = tagResponse
	TagRequest            = tagRequest
	ListTagsResponse      = listTagsResponse
	PostageCreateResponse = postageCreateResponse
	PostageStampResponse  = postageStampResponse
	PostageStampsResponse = postageStampsResponse
)

var (
	InvalidContentType  = errInvalidContentType
	InvalidRequest      = errInvalidRequest
	DirectoryStoreError = errDirectoryStore
)

var (
	ContentTypeTar    = contentTypeTar
	ContentTypeHeader = contentTypeHeader
)

var (
	ErrNoResolver           = errNoResolver
	ErrInvalidNameOrAddress = errInvalidNameOrAddress
)

var (
	FeedMetadataEntryOwner = feedMetadataEntryOwner
	FeedMetadataEntryTopic = feedMetadataEntryTopic
	FeedMetadataEntryType  = feedMetadataEntryType
)

func (s *Server) ResolveNameOrAddress(str string) (penguin.Address, error) {
	return s.resolveNameOrAddress(str)
}

func CalculateNumberOfChunks(contentLength int64, isEncrypted bool) int64 {
	return calculateNumberOfChunks(contentLength, isEncrypted)
}
