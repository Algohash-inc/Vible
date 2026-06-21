// models — list available models.
#include <iostream>

#include "vible.hpp"

int main() {
  vible::Global g;
  try {
    vible::Client c;
    for (const auto& id : c.models()) std::cout << id << "\n";
  } catch (const vible::Error& e) {
    std::cerr << "error: " << e.what() << "\n";
    return 1;
  }
  return 0;
}
