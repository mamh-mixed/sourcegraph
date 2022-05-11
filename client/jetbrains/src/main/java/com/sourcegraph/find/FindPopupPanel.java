package com.sourcegraph.find;

import com.intellij.openapi.Disposable;
import com.intellij.openapi.project.Project;
import com.intellij.openapi.ui.Splitter;
import com.intellij.ui.OnePixelSplitter;
import com.intellij.ui.PopupBorder;
import com.intellij.ui.components.JBPanel;
import com.intellij.ui.components.JBPanelWithEmptyText;
import com.intellij.ui.jcef.JBCefApp;
import com.intellij.util.ui.JBUI;
import com.sourcegraph.browser.JSToJavaBridgeRequestHandler;
import com.sourcegraph.browser.SourcegraphJBCefBrowser;
import org.cef.browser.CefBrowser;
import org.cef.browser.CefFrame;
import org.cef.handler.CefLoadHandler;
import org.cef.network.CefRequest;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

import javax.swing.*;
import java.awt.*;

/**
 * Inspired by <a href="https://sourcegraph.com/github.com/JetBrains/intellij-community/-/blob/platform/lang-impl/src/com/intellij/find/impl/FindPopupPanel.java">FindPopupPanel.java</a>
 */
public class FindPopupPanel extends JBPanel<FindPopupPanel> implements Disposable {
    private final SourcegraphJBCefBrowser browser;

    public FindPopupPanel(@NotNull Project project) {
        super(new BorderLayout());

        setPreferredSize(JBUI.size(1200, 800));
        setBorder(PopupBorder.Factory.create(true, true));
        setFocusCycleRoot(true);

        // Create splitter
        Splitter splitter = new OnePixelSplitter(true, 0.5f, 0.1f, 0.9f);
        add(splitter, BorderLayout.CENTER);

        PreviewPanel previewPanel = new PreviewPanel(project);

        // Create browser
        JBPanelWithEmptyText jcefPanel = new JBPanelWithEmptyText(new BorderLayout()).withEmptyText("Unfortunately, the browser is not available on your system. Try running the IDE with the default OpenJDK.");
        browser = JBCefApp.isSupported() ? new SourcegraphJBCefBrowser(new JSToJavaBridgeRequestHandler(project, previewPanel)) : null;
        if (browser != null) {
            jcefPanel.add(browser.getComponent(), BorderLayout.CENTER);
        }

        // Create top part
        JBPanelWithEmptyText jcefLoadingPanel = new JBPanelWithEmptyText(new BorderLayout());
        //noinspection DialogTitleCapitalization
        jcefLoadingPanel.getEmptyText().setText("Loading Sourcegraph...");
        //jcefLoadingPanel.setOpaque(true);
        JPanel topPanel = new JPanel();
        topPanel.setLayout(new OverlayLayout(topPanel));
        topPanel.add(jcefLoadingPanel);

        topPanel.add(jcefPanel); // This goes behind the other one

        if (browser != null) {
            browser.getJBCefClient().addLoadHandler(new CefLoadHandler() {
                @Override
                public void onLoadingStateChange(CefBrowser cefBrowser, boolean isLoading, boolean canGoBack, boolean canGoForward) {
                }

                @Override
                public void onLoadStart(CefBrowser cefBrowser, CefFrame frame, CefRequest.TransitionType transitionType) {
                }

                @Override
                public void onLoadEnd(CefBrowser cefBrowser, CefFrame frame, int httpStatusCode) {
                    //jcefLoadingPanel.setVisible(false);
                }

                @Override
                public void onLoadError(CefBrowser cefBrowser, CefFrame frame, ErrorCode errorCode, String errorText, String failedUrl) {
                }
            }, browser.getCefBrowser());
        }

        splitter.setFirstComponent(topPanel);
        splitter.setSecondComponent(previewPanel);
    }

    @Nullable
    public SourcegraphJBCefBrowser getBrowser() {
        return browser;
    }

    @Override
    public void dispose() {
        if (browser != null) {
            browser.dispose();
        }
    }
}
