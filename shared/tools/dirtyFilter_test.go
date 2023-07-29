package tools

import "testing"

func TestFilter(t *testing.T) {

    // 测试用例
    testCases := []struct {
        input    string
        expected string
    }{
        {"Hello world", "Hello world"},
        {"This is a test message", "This is a test message"},
        {"傻逼", "This is a ** message"},
        {"This is a shit message", "This is a **** message"},
        {"This is a badword message", "This is a ******* message"},
    }

    // 遍历测试用例
    for _, tc := range testCases {
        // 过滤敏感词
        output := Filter(tc.input)

        // 检查输出是否符合预期
        if output != tc.expected {
            t.Errorf("Filter(%q) = %q; expected %q", tc.input, output, tc.expected)
        }
    }
}
