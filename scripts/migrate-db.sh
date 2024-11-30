#!/bin/bash
goose postgres "postgres://root:password@localhost:5432/app" up -dir migrations