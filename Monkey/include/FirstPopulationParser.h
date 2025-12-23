#ifndef RECURSIVEDESCENT_PARSER_H
#define RECURSIVEDESCENT_PARSER_H
#include <iostream>
#include <utility>

#include "Lexer.h"

class FirstPopulationParser
{
public:
    explicit FirstPopulationParser(Lexer lexer)
        : m_lexer(std::move(lexer))
    {
        Advance();
    }

    bool Parse()
    {
        return Rule1() && m_currentToken.m_type == TokenType::T_EOF;
    }

private:
    Lexer m_lexer;
    Token m_currentToken;

    void Advance()
    {
        if (m_currentToken.m_type != TokenType::T_EOF)
        {
            std::cout << "Token: " << m_currentToken.m_lexeme << std::endl;
        }
        m_currentToken = m_lexer.GetNextToken();
    }

    bool Match(const TokenType type)
    {
        if (m_currentToken.m_type == type)
        {
            Advance();
            return true;
        }
        return false;
    }

    bool Rule1()
    {
        if (!Rule2())
        {
            return false;
        }
        return Rule1Prime();
    }

    bool Rule1Prime()
    {
        if (!Match(TokenType::T_AU))
        {
            return true;
        }
        if (!Rule2())
        {
            return false;
        }
        return Rule1Prime();
    }

    bool Rule2()
    {
        if (!Rule3())
        {
            return false;
        }
        return Rule2Prime();
    }

    bool Rule2Prime()
    {
        if (!Match(TokenType::T_KU))
        {
            return true;
        }
        if (!Rule3())
        {
            return false;
        }
        return Rule2Prime();
    }

    bool Rule3()
    {
        if (Match(TokenType::T_UH_TI))
        {
            return true;
        }
        if (Match(TokenType::T_HO))
        {
            return Rule3();
        }

        if (Match(TokenType::T_NU))
        {
            if (!Rule1())
            {
                return false;
            }
            return Match(TokenType::T_I_NU);
        }

        return false;
    }
};

#endif // RECURSIVEDESCENT_PARSER_H
