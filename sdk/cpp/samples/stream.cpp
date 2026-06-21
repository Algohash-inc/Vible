// stream — streaming chat completion; prints each token as it arrives.
//   ./stream "Write a haiku about the sea"
#include <iostream>

#include "vible.hpp"

int main(int argc, char** argv) {
  vible::Global g;
  std::string prompt = argc > 1 ? argv[1] : "Count from one to five.";
  try {
    vible::Client c;
    c.chat_stream({{"user", prompt}}, [](const std::string& token) {
      std::cout << token << std::flush;  // print incrementally
    });
    std::cout << "\n";
  } catch (const vible::Error& e) {
    std::cerr << "\nerror: " << e.what() << "\n";
    return 1;
  }
  return 0;
}
