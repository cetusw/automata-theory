#ifndef TOKEN_H
#define TOKEN_H

#include <string>
#include <string_view>

enum class TokenType {
    Keyword,
    Identifier,
    Integer,
    Float,
    StringLiteral,
    Operator,
    Punctuation,
    Unknown,
    EndOfFile
};

struct Token {
    TokenType type;
    std::string text;
    size_t offset;

    [[nodiscard]] std::string TypeToString() const {
        switch (type) {
            case TokenType::Keyword:       return "Keyword";
            case TokenType::Identifier:    return "Identifier";
            case TokenType::Integer:       return "Integer";
            case TokenType::Float:         return "Float";
            case TokenType::StringLiteral: return "StringLiteral";
            case TokenType::Operator:      return "Operator";
            case TokenType::Punctuation:   return "Punctuation";
            case TokenType::Unknown:       return "Unknown";
            case TokenType::EndOfFile:     return "EOF";
            default:                       return "Error";
        }
    }
};

#endif // TOKEN_H