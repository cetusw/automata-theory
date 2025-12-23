#ifndef LEXER_H
#define LEXER_H

#include "Token.h"
#include <string>
#include <vector>
#include <regex>
#include <unordered_set>

class Lexer {
public:
    explicit Lexer(std::string source);
    std::vector<Token> Tokenize();

private:
    std::string sourceCode;

    struct LexerRule {
        std::regex pattern;
        TokenType type;
    };

    std::vector<LexerRule> rules;
    static const std::unordered_set<std::string> keywords;

    void InitializeRules();
};

#endif // LEXER_H