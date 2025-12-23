#include <iostream>
#include <vector>
#include <iomanip>
#include "Lexer.h"

int main() {
    std::string source = R"(
        int main() {
            int a = 100;
            float b = 45.5;
            // Это комментарий
            if (a >= 100) {
                return "High";
            }
            return "Low";
        }
        int price = 100$; // $ вызовет ошибку (Unknown)
    )";

    std::cout << "Analyzing Source Code:\n" << source << "\n\n";
    std::cout << "---------------------------------------------------\n";

    Lexer lexer(source);
    std::vector<Token> tokens = lexer.Tokenize();

    for (const auto& token : tokens) {
        if (token.type == TokenType::EndOfFile) break;

        std::cout << std::left << std::setw(15) << token.TypeToString()
                  << " | " << token.text;

        if (token.type == TokenType::Unknown) {
            std::cout << "  <-- [LEXICAL ERROR]";
        }

        std::cout << std::endl;
    }

    return 0;
}