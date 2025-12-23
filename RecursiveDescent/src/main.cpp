#include <fstream>
#include <iostream>
#include "Lexer.h"
#include "Parser.h"

int main(const int argc, char* argv[])
{
    if (argc < 2)
    {
        std::cerr << "Usage: " << argv[0] << " <input_file>" << std::endl;
        return 1;
    }

    std::ifstream file(argv[1]);
    if (!file.is_open())
    {
        std::cerr << "Could not open file" << std::endl;
        return 1;
    }

    Lexer lexer(file);
    Parser parser(lexer);
    parser.Parse();

    std::cout << "Success: Program parsed correctly." << std::endl;
    return 0;
}
