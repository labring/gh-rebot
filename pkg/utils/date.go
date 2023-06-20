/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import "time"

func startOfWeek(t time.Time) time.Time {
	daysFromMonday := (int)(t.Weekday() - time.Monday)
	if daysFromMonday < 0 {
		daysFromMonday = 6
	}

	return t.AddDate(0, 0, -daysFromMonday).Truncate(24 * time.Hour)
}

func endOfWeek(t time.Time) time.Time {
	daysToSunday := (int)(7 - t.Weekday())
	if daysToSunday == 7 {
		daysToSunday = 0
	}

	return t.AddDate(0, 0, daysToSunday).Truncate(24 * time.Hour)
}

func FormatDay(t time.Time) string {
	return t.Format("2006/1/2")
}
func FormatWeek(t time.Time) (string, string) {
	start := startOfWeek(t)
	end := endOfWeek(t)

	return FormatDay(start), FormatDay(end)
}
