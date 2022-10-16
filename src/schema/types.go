// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema

// @min_length 1
type NonEmptyString = string

// @min_items 1
// @unique_items
type UniqueNonEmptyStringList = []NonEmptyString
