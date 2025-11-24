import "./Navigation.css";

export type PageType = "basic" | "gfm" | "code";

interface NavigationProps {
  currentPage: PageType;
  onPageChange: (page: PageType) => void;
}

interface NavItem {
  id: PageType;
  label: string;
}

const navItems: NavItem[] = [
  { id: "basic", label: "基本Markdown" },
  { id: "gfm", label: "GFM機能" },
  { id: "code", label: "コードハイライト" },
];

/**
 * ページ間のナビゲーションUIコンポーネント
 * シンプルなタブ型ナビゲーションを提供
 */
export function Navigation({ currentPage, onPageChange }: NavigationProps) {
  return (
    <nav className="navigation">
      <div className="nav-container">
        {navItems.map((item) => (
          <button
            key={item.id}
            className={`nav-button ${currentPage === item.id ? "active" : ""}`}
            onClick={() => onPageChange(item.id)}
          >
            {item.label}
          </button>
        ))}
      </div>
    </nav>
  );
}
