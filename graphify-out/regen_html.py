import json
from graphify.build import build_from_json
from graphify.export import to_html
from pathlib import Path

extraction = json.loads(Path('graphify-out/graph.json').read_text())
analysis = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())

G = build_from_json(extraction)
communities = {int(k): v for k, v in analysis['communities'].items()}
labels = {0: "Project Setup & Trading", 1: "SDK Core & Protocols"}

to_html(G, communities, 'graphify-out/graph.html', community_labels=labels)
print('HTML updated: ' + str(G.number_of_nodes()) + ' nodes')