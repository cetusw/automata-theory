#include <fstream>
#include <iostream>
#include "FirstPopulationParser.h"
#include "SecondPopulationParser.h"

int main(const int argc, char *argv[])
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


    Lexer lexer1(file);
    FirstPopulationParser p1(lexer1);
    Lexer lexer2(file);
    SecondPopulationParser p2;


    bool isPop1 = p1.Parse();
    bool isPop2 = p2.Parse(lexer2);


    std::cout << "Результат анализа:\n";

    if (isPop1)
    {
        std::cout << ">> Это ПЕРВАЯ популяция!\n";
    } else if (isPop2)
    {
        std::cout << ">> Это ВТОРАЯ популяция!\n";
    } else
    {
        std::cout << ">> Неизвестный вид (не подходит ни под одну грамматику).\n";
    }

    return 0;
}
