# go.env

```bash
echo "${var}"
echo "Substitute the value of var."

echo "${var:-word}"
echo "If var is null or unset, word is substituted for var. The value of var does not change."

echo "${var:=word}"
echo "If var is null or unset, var is set to the value of word."

echo "${var:?message}"
echo "If var is null or unset, message is printed to standard error. This checks that variables are set correctly."

echo "${var:+word}"
echo "If var is set, word is substituted for var. The value of var does not change."
```
