/*

Length-Prefix Value Encoding

Package lpve defines a variable-length encoding for binary octet strings up to
256 pebibytes in size. The encoded value itself is guaranteed not to exceed 72
octets. Values are self-parsing with deterministic lengths in every case.

Octet strings that are 64 octets or less in length are stored directly in the
encoded value. These values are said to be "inline".

Octet strings that are greater than 64 octets in length are represented by a
sakura encoded merkle tree root of the octet string. These values are said to
be "referenced". The sakura tree hashing mode, which specifies the actual
hash function to use, is implementation dependent and is not indicated in the
encoding.

Both inline and referenced values encode the length of the data that the
value represents. This has the following benefits:

1. References do not have to be resolved before the length of data they
  refer to can be determined.
2. Compression could be applied automatically for octet strings greater than
   some length.

Values are designed so that they naturally sort according to the length of
data that they represent (assuming big-endian interpretation with the leading
bit being the most significant).

Encoding overhead is kept proportional to octet length, so that small values
have less encoding overhead.

	X = value
	L = len(value) - 1
	J = len(value) - 65
	K = len(value) - 73
	P = len(len(value) - 65 >> 3)
	H = hash(value)
	M = merkle(value)

	Empty value        (1 byte):     00000000
	0-6 bit value      (1 byte):     01XXXXXX
	7-512 bit value    (2-65 bytes): 10LLLLLL XXXXXXXX [...]
	65-4096 byte value (2-72 bytes): 11000JJJ [JJJJJJJJ] [...] HHHHHHHH [...]
	65-4096 byte value (2-72 bytes): 11PPPKKK [KKKKKKKK] [...] HHHHHHHH [...]
	> 4096 byte value  (2-72 bytes): 11PPPKKK [KKKKKKKK] [...] MMMMMMMM [...]

Here are some examples of the length encoding for values longer than 64 bytes:

	11PPPJJJ JJJJJJJJ
	-------- --------
	11000000                                     =   65 +         0 +         0 +        0 =        65 bytes
	11000001                                     =   65 +         0 +         0 +        1 =        66 bytes
	11000010                                     =   65 +         0 +         0 +        2 =        67 bytes
	11000111                                     =   65 +         0 +         0 +        7 =        72 bytes
	11001000 00000000                            =   65 +         8 +         0 +        0 =        73 bytes
	11001000 00000001                            =   65 +         8 +         0 +        1 =        74 bytes
	11001000 10000000                            =   65 +         8 +         0 +      128 =       201 bytes
	11001000 11111111                            =   65 +         8 +         0 +      255 =       328 bytes
	11001001 00000000                            =   65 +         8 +       256 +        0 =       329 bytes
	11001001 11111111                            =   65 +         8 +       256 +      255 =       584 bytes
	11001010 00000000                            =   65 +         8 +       512 +        0 =       585 bytes
	11001010 11111111                            =   65 +         8 +       512 +      255 =       840 bytes
	11001011 00000000                            =   65 +         8 +       768 +        0 =       841 bytes
	11001011 11111111                            =   65 +         8 +       768 +      255 =      1096 bytes
	11001100 00000000                            =   65 +         8 +      1024 +        0 =      1097 bytes
	11001111 11111111                            =   65 +         8 +      1792 +      255 =      2120 bytes
	11010000 00000000 00000000                   =   65 +      2056 +         0 +        0 =      2121 bytes
	11010111 11111111 11111111                   =   65 +      2056 +    458752 +    65535 =    526408 bytes
	11011000 00000000 00000000 00000000          =   65 +    526344 +         0 +        0 =    526409 bytes
	11011111 11111111 11111111 11111111          =   65 +    526344 + 117440512 + 16777215 = 134744136 bytes
	11100000 00000000 00000000 00000000 00000000 =   65 + 134744072 +         0 +        0 = 134744137 bytes

*/
package lpve
