// call — place a cloud phone call to the workspace owner's device.
//   ./call "Your build just finished" ["title"]
// Needs a CALL key: set VIBLE_CALL_KEY, or run locally (resolved from the app).
#include <iostream>

#include "vible.hpp"

int main(int argc, char** argv) {
  vible::Global g;
  std::string prompt =
      argc > 1 ? argv[1] : "This is a test call from the Vible C++ SDK.";
  std::string title = argc > 2 ? argv[2] : "";
  try {
    vible::Client c;
    vible::CallResult r = c.call(prompt, title);
    std::cout << "call_id:   " << r.call_id << "\n";
    std::cout << "ringing:   " << (r.ringing ? "yes" : "no") << "\n";
    std::cout << "delivered: " << (r.delivered ? "yes" : "no") << "\n";
  } catch (const vible::Error& e) {
    std::cerr << "error: " << e.what() << "\n";
    return 1;
  }
  return 0;
}
