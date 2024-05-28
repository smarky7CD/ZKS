# ZKS

We provide an implementation of Zero-Knowledge Sets (ZKS) first introduced by Micali, Rabin, and Kilian [[MRK03](https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Zero%20Knowledge/Zero-Knowledge_Sets.pdf)]. We base our construction off the one described in *Mercurial Commitments with Applications to Zero-Knowledge Sets* by Chase et al. [[CHLMR05](https://cs.brown.edu/~mchase/papers/merc.pdf)].

## API

### Enum Sets

The data stored in a ZKS must be stored as an enumerated set. That is, items need to be stored in a map that maps integers (up to some max value) to a boolean value indicating set-membership.

For instance one can create a set with maximally 10 values as such:

```go
	values := map[uint64]bool{
		0:  true,
		1:  false,
		2:  false,
		3:  true,
		4:  true,
		5:  true,
		6:  false,
		7:  false,
		8:  false,
		9:  false,
		10: true,
    }

    set := NewEnumSet(values, 10)
```

One can also create this set dynamically by creating an empty EnumSet and using the `Add` and `Remove` functions. The `In` function answers a set-membership query. If the queried value is not explicitly stored in the map or is beyond the maximum value then `false` is returned.

### ZKS


- `Gen()` generates the public parameters for a ZKS (a selection of a verification key and a PRF).
- `Rep(pp,es)` takes as input the public parameters and an enumerated set. It outputs the ZKS representation and commitment to this representation. 
- `Qry(pp,repr,x)` takes as input the public parameters, the ZKS representation, and the element `x` being queried. It outputs the set-membership response and a proof to this response in a single `answer` struct. 
- `Vfy(pp,com,x,answer)` takes as input the public parameters, the commitment to the ZKS representation, the element `x` being queried, the answer/proof struct to a query on `x`. It outputs a boolean value indicating if the answer is valid.

## Installing and Using

To install the ZKS package:

```shell
go get -u github.com/smarky7CD/ZKS
```

The package can be imported and used as such:

```go
import(
    "fmt"
    "github.com/smarky7cd/ZKS"
)

func ZKSexample() {

	values := map[uint64]bool{
		0:  true,
		1:  false,
		2:  false,
		3:  true,
		4:  true,
		5:  true,
		6:  false,
		7:  false,
		8:  false,
		9:  false,
		10: true,
		11: true,
		12: true,
		13: true,
		14: true,
		15: true,
	}

	set := NewEnumSet(values, 16)

	pp := Gen()

	repr, com := Rep(pp, set)

	for i := uint64(0); i < 16; i++ {
		a := Qry(pp, repr, i)
        fmt.Println("Value ",i, " is: ", a.answer)
		v := Vfy(pp, com, i, a)
		fmt.Println("Answer is verified: ", v)
	}

}

```

## Tests 

We provide a number of correctness tests and a comprehensive performance evaluation.

They can be ran by using:

```shell
go test -v -timeout 5m
```

The performance tests will create a `.csv` file reporting the mean time and the variance for all operations for various set and universe sizes over 10 trials for each parameter set. The performance test may take a bit of time to run.

## References

- [MRK03]: Micali, Silvio, Michael Rabin, and Joe Kilian. "Zero-knowledge sets." 44th Annual IEEE Symposium on Foundations of Computer Science, 2003. Proceedings.. IEEE, 2003.
- [CHLMR05]: Chase, Melissa, et al. "Mercurial commitments with applications to zero-knowledge sets." Advances in Cryptology–EUROCRYPT 2005: 24th Annual International Conference on the Theory and Applications of Cryptographic Techniques, Aarhus, Denmark, May 22-26, 2005. Proceedings 24. Springer Berlin Heidelberg, 2005.
