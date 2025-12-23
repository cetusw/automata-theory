#ifndef RECURSIVEDESCENT_TOKEN_H
#define RECURSIVEDESCENT_TOKEN_H
#include <string>

enum class TokenType
{
    T_AU,
    T_KU,
    T_UH_TI,
    T_HO,
    T_NU,
    T_I_NU,
    T_OI,
    T_AI,
};

struct Token
{
    TokenType m_type;
    std::string m_lexeme;
};

#endif //RECURSIVEDESCENT_TOKEN_H