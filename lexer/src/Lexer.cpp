#include "Lexer.h"
#include <iostream>

const std::unordered_set<std::string> Lexer::keywords = {
    "int", "float", "return", "if", "else", "while", "for", "void", "class", "struct"
};

Lexer::Lexer(std::string source) : sourceCode(std::move(source)) {
    InitializeRules();
}

void Lexer::InitializeRules() {
    rules.push_back({std::regex(R"(^\d+\.\d+)"), TokenType::Float});
    rules.push_back({std::regex(R"(^\d+)"), TokenType::Integer});
    rules.push_back({std::regex(R"(^"[^"]*")"), TokenType::StringLiteral});
    rules.push_back({std::regex(R"(^[a-zA-Z_][a-zA-Z0-9_]*)"), TokenType::Identifier});
    rules.push_back({std::regex(R"(^(==|!=|<=|>=|&&|\|\|))"), TokenType::Operator});
    rules.push_back({std::regex(R"(^[+\-*/=<>!])"), TokenType::Operator});
    rules.push_back({std::regex(R"(^[;,\(\)\{\}\[\]])"), TokenType::Punctuation});
}

std::vector<Token> Lexer::Tokenize() {
    std::vector<Token> tokens;

    std::string::const_iterator currentIt = sourceCode.begin();
    const std::string::const_iterator endIt = sourceCode.end();
    size_t offset = 0;

    const std::regex whitespacePattern(R"(^\s+)");
    const std::regex commentPattern(R"(^//.*)");

    while (currentIt != endIt) {
        std::smatch match;

        if (std::regex_search(currentIt, endIt, match, whitespacePattern, std::regex_constants::match_continuous)) {
            const size_t length = match.length();
            currentIt += length;
            offset += length;
            continue;
        }

        if (std::regex_search(currentIt, endIt, match, commentPattern, std::regex_constants::match_continuous)) {
            size_t length = match.length();
            currentIt += length;
            offset += length;
            continue;
        }

        bool matched = false;
        for (const auto& rule : rules) {
            if (std::regex_search(currentIt, endIt, match, rule.pattern, std::regex_constants::match_continuous)) {
                std::string text = match.str();
                TokenType type = rule.type;

                if (type == TokenType::Identifier) {
                    if (keywords.contains(text)) {
                        type = TokenType::Keyword;
                    }
                }

                tokens.push_back({type, text, offset});

                const size_t length = match.length();
                currentIt += length;
                offset += length;
                matched = true;
                break;
            }
        }

        if (!matched) {
            tokens.push_back({TokenType::Unknown, std::string(1, *currentIt), offset});
            ++currentIt;
            ++offset;
        }
    }

    tokens.push_back({TokenType::EndOfFile, "", offset});
    return tokens;
}