#!/bin/bash
go test -v ../internal/utils
go test -v -cover ../internal/models
go test -v -cover ../internal/api
go test -cover