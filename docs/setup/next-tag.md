# Printing the next Tag Only

If all you need is the next calculated tag, Uplift can print this to `stdout` for you without making any changes to your repository. Useful if you want to use Uplift alongside other tools in your CI.

```sh
NEXT_TAG=$(uplift tag --next --silent)
```
