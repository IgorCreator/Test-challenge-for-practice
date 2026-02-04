export type TabKey = "store" | "add" | "history";

type TabsProps = {
  active: TabKey;
  onChange: (tab: TabKey) => void;
};

export default function Tabs({ active, onChange }: TabsProps) {
  return (
    <nav className="tabs">
      <button
        className={active === "store" ? "tab active" : "tab"}
        onClick={() => onChange("store")}
      >
        Purchase items
      </button>
      <button
        className={active === "add" ? "tab active" : "tab"}
        onClick={() => onChange("add")}
      >
        Add item
      </button>
      <button
        className={active === "history" ? "tab active" : "tab"}
        onClick={() => onChange("history")}
      >
        History
      </button>
    </nav>
  );
}
