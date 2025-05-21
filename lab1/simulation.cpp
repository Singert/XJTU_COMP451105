#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <map>
#include <set>
#include <queue>
using namespace std;

set<string> alphabet;
set<string> states;
string startState;
set<string> acceptStates;
map<pair<string, string>, string> transition;

void loadDFA(const string& filename) {
    ifstream file(filename);
    string line;

    while (getline(file, line)) {
        istringstream iss(line);
        string key;
        iss >> key;
        if (key == "alphabet:") {
            string sym;
            while (iss >> sym) alphabet.insert(sym);
        } else if (key == "states:") {
            string state;
            while (iss >> state) states.insert(state);
        } else if (key == "start:") {
            iss >> startState;
        } else if (key == "accept:") {
            string acc;
            while (iss >> acc) acceptStates.insert(acc);
        } else if (key == "transitions:") {
            break; // 后面是转移函数
        }
    }

    while (getline(file, line)) {
        istringstream iss(line);
        string from, symbol, to;
        iss >> from >> symbol >> to;
        transition[{from, symbol}] = to;
    }
}

bool checkDFAValid() {
    if (states.find(startState) == states.end()) {
        cout << "起始状态不在状态集中！\n";
        return false;
    }
    for (const string& s : acceptStates) {
        if (states.find(s) == states.end()) {
            cout << "接受状态" << s << "不在状态集中！\n";
            return false;
        }
    }
    return true;
}

bool runDFA(const string& input) {
    string current = startState;
    for (char c : input) {
        string sym(1, c);
        if (alphabet.find(sym) == alphabet.end()) return false;
        auto it = transition.find({current, sym});
        if (it == transition.end()) return false;
        current = it->second;
    }
    return acceptStates.find(current) != acceptStates.end();
}

void generateAllStrings(int N) {
    queue<pair<string, string>> q; // {current_string, current_state}
    q.push({"", startState});
    while (!q.empty()) {
        auto [s, state] = q.front(); q.pop();
        if (s.length() > N) continue;
        if (acceptStates.count(state)) cout << s << " [ACCEPTED]\n";

        for (const string& sym : alphabet) {
            auto it = transition.find({state, sym});
            if (it != transition.end()) {
                q.push({s + sym, it->second});
            }
        }
    }
}

int main() {
    loadDFA("../dfa.txt");

    if (!checkDFAValid()) return 1;

    cout << "\n生成所有长度 <= 3 的合法字符串：\n";
    generateAllStrings(3);

    string testInput;
    cout << "\n请输入字符串判断是否被DFA接受（输入end结束）：\n";
    while (cin >> testInput && testInput != "end") {
        if (runDFA(testInput)) {
            cout << "字符串被接受！\n";
        } else {
            cout << "字符串不被接受。\n";
        }
    }

    return 0;
}

