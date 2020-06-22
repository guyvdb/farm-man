package sequence

// A sequence is a string Id that has a prefix and a monotonically increasing number such as INV000082
// The number portion of a sequence is six digits long. I.e. XXX000001, XXX000002, etc

type Sequence string
