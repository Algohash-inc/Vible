// transcribe — read a wav file and print the transcript.
//   ./transcribe [path]   (defaults to out.wav — run ./speech first)
#include <fstream>
#include <iostream>
#include <iterator>

#include "vible.hpp"

int main(int argc, char** argv) {
  vible::Global g;
  std::string path = argc > 1 ? argv[1] : "out.wav";
  std::ifstream f(path, std::ios::binary);
  if (!f) {
    std::cerr << "cannot open " << path << " (run ./speech first?)\n";
    return 1;
  }
  std::vector<uint8_t> audio((std::istreambuf_iterator<char>(f)),
                             std::istreambuf_iterator<char>());
  try {
    vible::Client c;
    std::cout << c.transcribe(audio) << "\n";
  } catch (const vible::Error& e) {
    std::cerr << "error: " << e.what() << "\n";
    return 1;
  }
  return 0;
}
