import { useState } from "react";
import { Navigation, type PageType } from "./components/Navigation";
import { BasicMarkdownDemo } from "./pages/BasicMarkdownDemo";
import { GfmDemo } from "./pages/GfmDemo";
import { CodeHighlightDemo } from "./pages/CodeHighlightDemo";
import "./App.css";

function App() {
  const [currentPage, setCurrentPage] = useState<PageType>("basic");

  const renderPage = () => {
    switch (currentPage) {
      case "basic":
        return <BasicMarkdownDemo />;
      case "gfm":
        return <GfmDemo />;
      case "code":
        return <CodeHighlightDemo />;
      default:
        return <BasicMarkdownDemo />;
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>React Markdown 検証サイト</h1>
        <p className="subtitle">
          react-markdown + remark-gfm + react-syntax-highlighter
        </p>
      </header>

      <Navigation currentPage={currentPage} onPageChange={setCurrentPage} />

      <main className="app-main">{renderPage()}</main>

      <footer className="app-footer">
        <p>Built with React + Vite + TypeScript</p>
      </footer>
    </div>
  );
}

export default App;
