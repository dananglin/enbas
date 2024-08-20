# Tips and Tricks

## The time duration value

The time duration value is a custom [flag value](https://pkg.go.dev/flag#Value) that converts a string input into a duration of time.
A typical string input would be in the form of something like `"3 days, 12 hours and 39 minutes"`.
The value can convert units in days, hours, minutes and seconds.

To ensure that your string input is converted correctly there are simple rules to follow.

- The input must be wrapped in quotes.
- Use `day` or `days` to convert the number of days.
- Use `hour` or `hours` to convert the number of hours.
- Use `minute` or `minutes` to convert the number of minutes.
- Use `second` or `seconds` to convert the number of seconds.
- There must be at least one space between the number and the unit of time.<br>
  E.g. `"7 days"` is valid, but `"7days"` is invalid.

### Example valid string inputs

- `"3 days"`
- `"6 hours, 45 minutes and 1 second"`
- `"1 day, 15 hours 31 minutes and 12 seconds"`
- `"(7 days) (1 hour) (21 minutes) (35 seconds)"`
