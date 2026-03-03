apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: crypto-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/huzaifa678/Continious-Delivery.git
    targetRevision: main
    path: eks-chart
    helm:
      valueFiles:
        - values.yaml
      parameters:
        - name: environment
          value: "${environment}"
  destination:
    server: https://kubernetes.default.svc
    namespace: my-app
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
