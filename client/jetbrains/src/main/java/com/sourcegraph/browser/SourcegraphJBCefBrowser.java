package com.sourcegraph.browser;

import com.intellij.openapi.project.Project;
import com.intellij.openapi.util.Disposer;
import com.intellij.ui.jcef.JBCefBrowser;
import com.sourcegraph.config.ThemeUtil;
import org.cef.CefApp;
import org.jetbrains.annotations.NotNull;

public class SourcegraphJBCefBrowser extends JBCefBrowser {
    private final JSToJavaBridge jsToJavaBridge;

    public SourcegraphJBCefBrowser(@NotNull Project project) {
        super("http://sourcegraph/html/index.html");
        // Create and set up JCEF browser
        CefApp.getInstance().registerSchemeHandlerFactory("http", "sourcegraph", new HttpSchemeHandlerFactory());
        this.setPageBackgroundColor(ThemeUtil.getPanelBackgroundColorHexString());

        // Create bridge, set up handlers, then run init function
        String initJSCode = "window.initializeSourcegraph();";
        jsToJavaBridge = new JSToJavaBridge(this, new JSToJavaBridgeRequestHandler(project), initJSCode);
        Disposer.register(this, jsToJavaBridge);
    }

    public JSToJavaBridge getJsToJavaBridge() {
        return jsToJavaBridge;
    }

    public void focus() {
        this.getCefBrowser().setFocus(true);
    }
}
