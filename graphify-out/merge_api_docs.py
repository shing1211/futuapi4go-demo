import json
from pathlib import Path
from networkx.readwrite import json_graph

# Load existing graph
existing = json.loads(Path('graphify-out/graph.json').read_text())
G = json_graph.node_link_graph(existing, edges='links')

# Load new API docs
api = json.loads(Path('graphify-out/.graphify_api_docs.json').read_text())

# Add new nodes
for n in api['nodes']:
    if n['id'] not in G:
        G.add_node(n['id'], **n)

# Add new edges
for e in api['edges']:
    if not G.has_edge(e['source'], e['target']):
        G.add_edge(e['source'], e['target'], **e)

print(f'Updated graph: {G.number_of_nodes()} nodes, {G.number_of_edges()} edges')

# Rebuild communities and analysis
from graphify.cluster import cluster, score_all
from graphify.analyze import god_nodes, surprising_connections, suggest_questions
from graphify.export import to_json
from graphify.report import generate

detection = {'total_files': 97, 'total_words': 22946, 'needs_graph': True, 'warning': None}
communities = cluster(G)
cohesion = score_all(G, communities)
gods = god_nodes(G)
surprises = surprising_connections(G, communities)
labels = {0: "Project Setup & Trading", 1: "SDK Core & Protocols", 9: "K-Line Data", 10: "Error Handling"}
questions = suggest_questions(G, communities, labels)

to_json(G, communities, 'graphify-out/graph.json')

report = generate(G, communities, cohesion, labels, gods, surprises, detection, {'input': 17000, 'output': 6500}, '.', suggested_questions=questions)
Path('graphify-out/GRAPH_REPORT.md').write_text(report, encoding='utf-8')

analysis = {
    'communities': {str(k): v for k, v in communities.items()},
    'cohesion': {str(k): v for k, v in cohesion.items()},
    'gods': gods,
    'surprises': surprises,
    'questions': questions,
}
Path('graphify-out/.graphify_analysis.json').write_text(json.dumps(analysis, indent=2))
print('Graph rebuilt with API docs')