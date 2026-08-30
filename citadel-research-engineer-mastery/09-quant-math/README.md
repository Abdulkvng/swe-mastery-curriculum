# Phase 09 — Quantitative Math Foundation

You do not need to become a quant researcher overnight. You do need fast, precise reasoning with probability, statistics, linear algebra, and numerical computation.

## Probability essentials

Know:

- conditional probability
- Bayes' rule
- independence
- expectation
- variance / covariance
- common distributions
- law of total expectation

### Expectation

For discrete X:

`E[X] = Σ x P(X=x)`

Linearity is powerful:

`E[aX+bY] = aE[X] + bE[Y]`

No independence required.

### Variance

`Var(X)=E[X²]-E[X]²`

For independent X,Y:

`Var(X+Y)=Var(X)+Var(Y)`

Without independence include covariance.

## Conditional probability

`P(A|B) = P(A∩B)/P(B)`

Bayes:

`P(A|B)=P(B|A)P(A)/P(B)`

Interview questions often test whether you condition on the right information.

## Statistics

Understand:

- sample mean and variance
- estimator bias
- confidence intervals conceptually
- correlation vs causation
- overfitting
- train/test leakage
- multiple testing intuition
- stationarity intuition

## Regression

Linear model:

`y = Xβ + ε`

Least squares minimizes squared residuals. Know what a coefficient means, why correlated features cause instability, and why good in-sample fit does not guarantee predictive value.

## Time-series intuition

Financial observations are ordered and dependence matters.

Know concepts:

- returns vs prices
- autocorrelation
- rolling statistics
- regime change
- look-ahead bias
- nonstationarity

## Numerical issues

### Floating point is finite

You must understand:

- rounding error
- catastrophic cancellation
- associativity does not hold exactly
- overflow/underflow
- `float` vs `double`

Example: summation order can change a floating-point result. Parallelizing a reduction can therefore change the last bits.

### Stable algorithms

If subtracting nearly equal large values destroys precision, reformulate when possible. Know Kahan summation conceptually for reducing accumulation error.

## Mental-math drills

Practice:

- expected value of simple games
- conditional probability from tables
- covariance signs
- percentage/basis-point changes
- weighted averages
- simple combinatorics

## Interview problems

1. Two independent fair dice: expected maximum?
2. Coin with unknown bias and observations: how does Bayesian reasoning change belief?
3. Why does correlation between a feature and future returns not prove a tradable signal?
4. Why is random train/test splitting dangerous for time series?
5. How could changing summation order alter a production model's output?

Use `quizzes/quant-quiz.md` for closed-book practice.
