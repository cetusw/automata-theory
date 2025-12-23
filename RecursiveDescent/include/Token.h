#ifndef RECURSIVEDESCENT_TOKEN_H
#define RECURSIVEDESCENT_TOKEN_H
#include <string>

enum class TokenType
{
    T_MAIN,
    T_END,
    T_BEGIN,
    T_VAR,
    T_CONST,
    T_INT,
    T_FLOAT,
    T_ID,
    T_NUMBER,
    T_COLON,
    T_SEMICOLON,
    T_COMMA,
    T_DOT,
    T_ASSIGN,
    T_EQUALS,
    T_PLUS,
    T_MULTIPLICATION,
    T_MINUS,
    T_LEFT_PARENTHESIS,
    T_RIGHT_PARENTHESIS,
    T_EOF,
    T_ERROR,
};

struct Token
{
    TokenType m_type;
    std::string m_lexeme;
};

#endif //RECURSIVEDESCENT_TOKEN_H