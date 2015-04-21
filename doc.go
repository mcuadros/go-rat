/*
rat is an extension to the classical tar archive, focused on allowing
constant-time random file access with linear memory consumption increase. tape
archive, was originally developed to write and read streamed sources, making
random access to the content very inefficient.

Based on the benchmarks, we found that rat is 4x to 60x times faster over SSD
and HDD than the classic tar file, when reading a single file from a tar
archive.

Note: Any tar file produced by rat is compatible with standard tar
implementation.
*/
package rat
