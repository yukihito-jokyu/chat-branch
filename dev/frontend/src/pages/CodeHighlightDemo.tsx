import { MarkdownRenderer } from "../components/MarkdownRenderer";

const codeHighlightContent = `# コードハイライト機能のデモ

このページでは、react-syntax-highlighterを使用したシンタックスハイライト機能を紹介します。

## JavaScript

\`\`\`javascript
// 関数の定義
function fibonacci(n) {
  if (n <= 1) return n;
  return fibonacci(n - 1) + fibonacci(n - 2);
}

// 配列操作
const numbers = [1, 2, 3, 4, 5];
const doubled = numbers.map(n => n * 2);
console.log(doubled); // [2, 4, 6, 8, 10]

// 非同期処理
async function fetchData(url) {
  try {
    const response = await fetch(url);
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error:', error);
  }
}
\`\`\`

## TypeScript

\`\`\`typescript
// インターフェースの定義
interface User {
  id: number;
  name: string;
  email: string;
  role: 'admin' | 'user' | 'guest';
}

// ジェネリック型の使用
function createArray<T>(length: number, value: T): T[] {
  return Array(length).fill(value);
}

// クラスの定義
class UserManager {
  private users: User[] = [];

  addUser(user: User): void {
    this.users.push(user);
  }

  findUserById(id: number): User | undefined {
    return this.users.find(user => user.id === id);
  }
}
\`\`\`

## Python

\`\`\`python
# リスト内包表記
squares = [x**2 for x in range(10)]

# デコレータ
def timer(func):
    def wrapper(*args, **kwargs):
        import time
        start = time.time()
        result = func(*args, **kwargs)
        end = time.time()
        print(f"{func.__name__} took {end - start:.2f} seconds")
        return result
    return wrapper

@timer
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

# クラスの定義
class Person:
    def __init__(self, name, age):
        self.name = name
        self.age = age
    
    def greet(self):
        return f"Hello, I'm {self.name} and I'm {self.age} years old."
\`\`\`

## JSON

\`\`\`json
{
  "name": "react-markdown-demo",
  "version": "1.0.0",
  "dependencies": {
    "react": "^19.2.0",
    "react-markdown": "^10.1.0",
    "react-syntax-highlighter": "^16.1.0",
    "remark-gfm": "^4.0.1"
  },
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build"
  }
}
\`\`\`

## CSS

\`\`\`css
/* グラデーション背景 */
.gradient-bg {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  border-radius: 12px;
}

/* フレックスボックスレイアウト */
.container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

/* ホバーエフェクト */
.button {
  transition: all 0.3s ease;
  transform: translateY(0);
}

.button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
\`\`\`

## インラインコードとの比較

インラインコード: \`const greeting = "Hello";\` は文中に埋め込まれます。

一方、コードブロックは独立した表示になり、シンタックスハイライトと行番号が適用されます：

\`\`\`javascript
const greeting = "Hello";
console.log(greeting);
\`\`\`
`;

/**
 * コードハイライト機能のデモページ
 */
export function CodeHighlightDemo() {
  return (
    <div className="demo-page">
      <MarkdownRenderer
        content={codeHighlightContent}
        className="markdown-content"
      />
    </div>
  );
}
