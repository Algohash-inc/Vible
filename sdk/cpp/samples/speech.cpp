// speech — synthesize speech and write out.wav.
//   ./speech "Hello from Vible" [voice]
// voice: ""/alloy/echo/... → default; "clone:NAME" → local clone; other → speaker name.
#include <fstream>
#include <iostream>

#include "vible.hpp"

int main(int argc, char** argv) {
  vible::Global g;
  std::string text = argc > 1 ? argv[1] : "Hello from the Vible C++ SDK.";
  vible::SpeechOptions opt;
  if (argc > 2) opt.voice = argv[2];
  try {
    vible::Client c;
    std::vector<uint8_t> wav = c.speech(text, opt);
    if (wav.empty()) {
      std::cout << "spoken on device (no audio bytes returned)\n";
      return 0;
    }
    std::ofstream f("out.wav", std::ios::binary);
    f.write(reinterpret_cast<const char*>(wav.data()),
            static_cast<std::streamsize>(wav.size()));
    std::cout << "wrote out.wav (" << wav.size() << " bytes)\n";
  } catch (const vible::Error& e) {
    std::cerr << "error: " << e.what() << "\n";
    return 1;
  }
  return 0;
}
