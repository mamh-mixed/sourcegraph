export enum PageRoutes {
    Index = '/',
    Search = '/search',
    SearchConsole = '/search/console',
    SearchNotebook = '/search/notebook',
    SignIn = '/sign-in',
    SignUp = '/sign-up',
    UnlockAccount = '/unlock-account/:token',
    Welcome = '/welcome',
    Settings = '/settings',
    User = '/user',
    Organizations = '/organizations',
    SiteAdmin = '/site-admin',
    SiteAdminInit = '/site-admin/init',
    PasswordReset = '/password-reset',
    ApiConsole = '/api/console',
    UserArea = '/users/:username',
    Survey = '/survey/:score?',
    Extensions = '/extensions',
    Help = '/help',
    Debug = '/-/debug/*',
    NotebookCreate = '/notebooks/new',
    Notebook = '/notebooks/:id',
    Notebooks = '/notebooks',
    RepoContainer = '/:repoRevAndRest+',
    InstallGitHubAppSuccess = '/install-github-app-success',
    InstallGitHubAppSelectOrg = '/install-github-app-select-org',
}

export enum EnterprisePageRoutes {
    SubscriptionsNew = '/subscriptions/new',
    OldSubscriptionsNew = '/user/subscriptions/new',
    BatchChanges = '/batch-changes',
    Stats = '/stats',
    CodeMonitoring = '/code-monitoring',
    Insights = '/insights',
    Contexts = '/contexts',
    CreateContext = '/contexts/new',
    EditContext = '/contexts/:spec+/edit',
    Context = '/contexts/:spec+',
}
