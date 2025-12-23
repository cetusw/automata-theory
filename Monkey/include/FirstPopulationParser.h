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

    void Parse()
    {
        Rule1();
        if (m_currentToken.m_type != TokenType::T_EOF)
        {
            Error("Extra tokens after end of program");
        }
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
        if (m_currentToken.m_type == TokenType::T_ERROR)
        {
            Error("Lexical error: " + m_currentToken.m_lexeme);
        }
    }

    void Error(const std::string &message) const
    {
        std::cerr << "Error: " << message << " (Current: " << m_currentToken.m_lexeme << ")" << std::endl;
        exit(1);
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

    void Expect(const TokenType type, const std::string &message)
    {
        if (!Match(type))
        {
            Error(message);
        }
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

    void Declaration()
    {
        if (m_currentToken.m_type == TokenType::T_VAR)
        {
            Var();
        } else if (m_currentToken.m_type == TokenType::T_CONST || m_currentToken.m_type == TokenType::T_ID)
        {
            Consts();
        } else
        {
            Error("Expected declaration (var or const)");
        }
    }

    void Var()
    {
        Expect(TokenType::T_VAR, "Expected 'var'");
        IdList();
        Expect(TokenType::T_COLON, "Expected ':'");
        Type();
    }

    void IdList()
    {
        Expect(TokenType::T_ID, "Expected identifier");
        while (Match(TokenType::T_COMMA))
        {
            Expect(TokenType::T_ID, "Expected identifier after ','");
        }
    }

    void Type()
    {
        if (m_currentToken.m_type == TokenType::T_INT || m_currentToken.m_type == TokenType::T_FLOAT)
        {
            Advance();
        } else
        {
            Error("Expected type 'int' or 'float'");
        }
    }

    void Consts()
    {
        Expect(TokenType::T_CONST, "Expected 'const'");
        Const();
        while (m_currentToken.m_type == TokenType::T_SEMICOLON)
        {
            const Token next = m_lexer.PeekToken();
            if (next.m_type == TokenType::T_ID)
            {
                Advance();
                Const();
            } else
            {
                break;
            }
        }
    }

    void Const()
    {
        Expect(TokenType::T_ID, "Expected identifier in constant");
        Expect(TokenType::T_EQUALS, "Expected '='");
        Expression();
    }

    void Statements()
    {
        Statement();
        while (m_currentToken.m_type == TokenType::T_SEMICOLON)
        {
            const Token next = m_lexer.PeekToken();
            if (next.m_type == TokenType::T_ID)
            {
                Advance();
                Statement();
            } else
            {
                break;
            }
        }
    }

    void Statement()
    {
        Assign();
    }

    void Assign()
    {
        Expect(TokenType::T_ID, "Expected identifier in assignment");
        Expect(TokenType::T_ASSIGN, "Expected ':='");
        Expression();
    }

    void Expression()
    {
        T();
        while (Match(TokenType::T_PLUS))
        {
            T();
        }
    }

    void T()
    {
        F();
        while (Match(TokenType::T_MULTIPLICATION))
        {
            F();
        }
    }

    void F()
    {
        if (Match(TokenType::T_MINUS))
        {
            F();
        } else if (Match(TokenType::T_LEFT_PARENTHESIS))
        {
            Expression();
            Expect(TokenType::T_RIGHT_PARENTHESIS, "Expected ')'");
        } else if (!Match(TokenType::T_ID) && !Match(TokenType::T_NUMBER))
        {
            Error("Expected factor (id, number, '-' or '(')");
        }
    }
};

#endif // RECURSIVEDESCENT_PARSER_H
