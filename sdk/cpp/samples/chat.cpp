// chat — non-streaming chat completion.
//   ./chat "Tell me a joke about steel"
#include <iostream>

#include "vible.hpp"

int main(int argc, char** argv) {
  vible::Global g;  // curl_global_init / cleanup
  std::string prompt = argc > 1 ? argv[1] : "Say hello in one sentence.";
  try {
    vible::Client c;
    std::string reply = c.chat({{"user", prompt}});
    std::cout << reply << "\n";
  } catch (const vible::Error& e) {
    std::cerr << "error: " << e.what() << "\n";
    return 1;
  }
  return 0;
}
