const INIT_KEY = "kb-nav-default-state-initialized";

function getStoredFlag() {
  try {
    return window.sessionStorage.getItem(INIT_KEY) === "1";
  } catch {
    return false;
  }
}

function setStoredFlag() {
  try {
    window.sessionStorage.setItem(INIT_KEY, "1");
  } catch {
    // Ignore storage failures and fall back to best-effort behavior.
  }
}

function navLabelText(item) {
  const label = item?.querySelector(":scope > label .md-ellipsis");
  return label?.textContent?.trim() || "";
}

function setExpanded(item, expanded) {
  const toggle = item?.querySelector(":scope > input.md-nav__toggle");
  const nav = item?.querySelector(":scope > nav.md-nav");
  if (!toggle) {
    return false;
  }

  toggle.checked = expanded;
  if (nav) {
    nav.setAttribute("aria-expanded", expanded ? "true" : "false");
  }
  return true;
}

function applyKbNavDefaultState() {
  if (getStoredFlag()) {
    return;
  }

  const topLevelItems = document.querySelectorAll(
    ".md-nav--primary > .md-nav__list > .md-nav__item",
  );
  const sourcesItem = Array.from(topLevelItems).find(
    (item) => navLabelText(item) === "Sources",
  );
  if (!sourcesItem) {
    return;
  }

  const changed = setExpanded(sourcesItem, true);
  const yearItems = sourcesItem.querySelectorAll(
    ":scope > nav > .md-nav__list > .md-nav__item--nested",
  );
  yearItems.forEach((item) => {
    setExpanded(item, true);
  });

  if (changed) {
    setStoredFlag();
  }
}

function scheduleKbNavDefaultState() {
  window.requestAnimationFrame(applyKbNavDefaultState);
}

if (window.document$?.subscribe) {
  window.document$.subscribe(scheduleKbNavDefaultState);
} else if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", scheduleKbNavDefaultState, {
    once: true,
  });
} else {
  scheduleKbNavDefaultState();
}
